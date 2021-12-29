import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import Login from "./Login";
import Results from "./Results";
import Me from "./Me";
import Summary from "./Summary";

const App = () => (
  <Router>
    <Routes>
      <Route path="/">
        <Route index element={<Login />} />
        <Route path="me" element={<Me />} />
        <Route path="summary/:id" element={<Summary />} />
        <Route path="results/:id" element={<Results />} />
      </Route>
    </Routes>
  </Router>
);

export default App;
