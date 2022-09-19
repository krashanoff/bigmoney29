import { configureStore } from "@reduxjs/toolkit";

import authReducer from "./auth";
import assnReducer from "./assignments";

export default configureStore({
  reducer: {
    auth: authReducer,
    assignments: assnReducer,
  },
});
