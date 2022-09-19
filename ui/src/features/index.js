import { configureStore } from "@reduxjs/toolkit";

import authReducer from "./auth";
import assnReducer from "./assignments";

const loadState = () => {
  try {
    const state = localStorage.getItem("state");
    if (!state) {
      return undefined;
    }
    return JSON.parse(state);
  } catch (e) {
    return undefined;
  }
};

const saveState = (state) => {
  localStorage.setItem("state", JSON.stringify(state));
};

const persistedState = loadState();

const state = configureStore({
  reducer: {
    auth: authReducer,
    assignments: assnReducer,
  },
  preloadedState: persistedState,
});

state.subscribe(() => {
  saveState(state.getState());
});

export default state;
