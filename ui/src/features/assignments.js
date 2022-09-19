import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { AUTHORIZED_API_BASE } from "./constant";

/**
 * Update the state's current view of assignments in the system.
 */
const updateAssignments = createAsyncThunk(
  "api/updateAssignments",
  async (_noArgs, thunkAPI) => {
    const state = thunkAPI.getState();
    console.info("Updating assignments...");

    if (!state.auth.token) {
      thunkAPI.dispatch("User is not logged in");
      return;
    }

    try {
      thunkAPI.dispatch(assnSlice.actions.setFetch(true));
      const response = await fetch(`${AUTHORIZED_API_BASE}/assignments`, {
        headers: {
          Authorization: `Bearer ${state.auth.token}`,
        },
      });
      thunkAPI.dispatch(assnSlice.actions.setFetch(false));

      if (response.status !== 200) {
        throw new Error(`Server responded with code ${response.status}`);
      }

      const { assignments } = await response.json();
      return assignments;
    } catch (e) {
      thunkAPI.dispatch(assnSlice.actions.setFetch(false));
      thunkAPI.dispatch(assnSlice.actions.setError(e.message));
      return [];
    }
  }
);

const assnSlice = createSlice({
  name: "assignments",
  initialState: {
    loading: false,
    assignments: [],
    fetching: false,
    errorMessage: "",
  },
  reducers: {
    setError: (state, action) => {
      state.errorMessage = action.payload;
    },
    clearError: (state, actions) => {
      state.errorMessage = "";
    },
    setFetch: (state, action) => {
      state.fetching = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(updateAssignments.fulfilled, (state, action) => {
      state.fetching = false;
      state.assignments = action.payload;
    });
  },
});

export { updateAssignments };
export const { setError, setFetch } = assnSlice.actions;
export default assnSlice.reducer;
