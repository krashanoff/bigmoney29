import React from 'react';

const JWTContext = React.createContext({
  jwt: "",
  setJWT: (j) => {},

  assignments: [],
  updateAssignments: () => {},
});

export default JWTContext;
