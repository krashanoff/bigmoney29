import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import {
  AUTHORIZED_API_BASE,
  REQUEST_TIMEOUT,
  UNAUTHORIZED_API_BASE,
} from "./constant";

/**
 * Log into the service, returning a token for future authorized requests,
 * or setting the error message.
 *
 * Requires a form element with `username` and `password` fields.
 *
 * Sets the `auth.fetching` state variable while waiting, with a default timeout
 * of `DEFAULT_TIMEOUT` milliseconds. In all cases, this operation will terminate
 * with `auth.fetching === false`.
 */
const login = createAsyncThunk("auth/login", async (formElement, thunkAPI) => {
  const state = thunkAPI.getState();
  thunkAPI.dispatch(authSlice.actions.clearError());
  console.info("Logging in...");
  if (state.auth.fetching) {
    console.info("In progress already... Aborting.");
    return "";
  }

  if (!state.auth.token) {
    thunkAPI.dispatch(authSlice.actions.setFetch(true));

    const returnReset = (v) => {
      thunkAPI.dispatch(authSlice.actions.setFetch(false));
      return v;
    };

    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), REQUEST_TIMEOUT);

    try {
      const response = await fetch(`${UNAUTHORIZED_API_BASE}/login`, {
        method: "post",
        body: new FormData(formElement),
        signal: controller.signal,
      });
      clearTimeout(id);

      if (response.status === 403) {
        thunkAPI.dispatch(
          authSlice.actions.setError("Username or password is incorrect")
        );
        return returnReset("");
      }

      if (response.status !== 200) {
        thunkAPI.dispatch(
          authSlice.actions.setError(
            `Server responded with code ${response.status}`
          )
        );
        return returnReset("");
      }

      const { token } = await response.json();
      return returnReset(token);
    } catch (e) {
      thunkAPI.dispatch(authSlice.actions.setFetch(false));
      thunkAPI.dispatch(authSlice.actions.setError(e.message));
      clearTimeout(id);
      return returnReset("");
    }
  } else {
    return state.auth.token;
  }
});

/**
 * Log an user out of the system, invalidating their current session token.
 *
 * Like `login`, `auth.fetching` will be set to `true` while making the request,
 * then set to `false` after running. The request will terminate in the event
 * that `DEFAULT_TIMEOUT` milliseconds have passed without a response.
 */
const logout = createAsyncThunk("auth/logout", async (_unused, thunkAPI) => {
  const state = thunkAPI.getState();
  if (!state.auth.token) {
    thunkAPI.dispatch(authSlice.actions.setError("Not logged in!"));
    return;
  }

  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), 3000);

  try {
    const response = await fetch(`${AUTHORIZED_API_BASE}/logout`, {
      method: "post",
      body: JSON.stringify({
        token: state.auth.token,
      }),
      signal: controller.signal,
    });
    clearTimeout(id);
    if (response.status !== 200) {
      throw new Error("Couldn't log out!");
    }
  } catch (e) {
    thunkAPI.dispatch(authSlice.actions.setError(e.message));
    clearTimeout(id);
  }
});

const authSlice = createSlice({
  name: "auth",
  initialState: {
    token: "",
    fetching: false,
    errorMessage: "",
  },
  reducers: {
    setError: (state, action) => {
      state.errorMessage = action.payload;
    },
    clearError: (state, action) => {
      state.errorMessage = "";
    },
    setFetch: (state, action) => {
      state.fetching = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(login.fulfilled, (state, action) => {
      state.loading = false;
      state.token = action.payload;
    });
  },
});

export { login, logout };
export const { clearError } = authSlice.actions;
export default authSlice.reducer;
