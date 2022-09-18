import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import Login from "./Login";
import Results from "./Results";
import Me from "./Me";
import Assignment from "./Assignment";

const App = () => (
  <Router>
    <Routes>
      <Route path="/">
        <Route index element={<Login />} />
        <Route path="me" element={<Me />} />
        <Route path="assignment/:id" element={<Assignment />} />
        <Route path="results/:id" element={<Results />} />
      </Route>
    </Routes>
  </Router>
);

export default App;
