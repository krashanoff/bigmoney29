import { BrowserRouter as Router, Route, Routes } from "react-router-dom";

import Login from './Login';
import Results from "./Results";
import Me from './Me';
import Summary from "./Summary";
import JWTContext from "./jwt-context";
import { useState } from "react";

const App = () => {
  const [jwtState, setJwtState] = useState({
    jwt: "",
    setJWT: (newToken) => setJwtState((j) => { j.jwt = newToken; return j; }),
    updateAssignments: () => {
      fetch("http://localhost:8081/cash/assignments", {
        method: "get",
        mode: 'cors',
        headers: {
          "Authorization": `Bearer: ${jwtState.jwt}`,
        },
      })
      .then(r => r.text())
      .then(t => console.info);
    },
  });

  return (
    <JWTContext.Provider value={jwtState}>
      <Router>
        <Routes>
          <Route path="/">
            <Route index element={<Login />} />
            <Route path="me" element={<Me />} />
            <Route path="results/">
              <Route index element={<Summary />} />
              <Route path=":id" element={<Results />} />
            </Route>
          </Route>
        </Routes>
      </Router>
    </JWTContext.Provider>
  );
};

export default App;
