import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import Collapse from "react-bootstrap/Collapse";
import Alert from "react-bootstrap/Alert";

import { useDispatch, useSelector } from "react-redux";
import { clearError, login } from "./features/auth";

const Home = () => {
  const state = useSelector((state) => state.auth);
  const dispatch = useDispatch();

  return (
    <div
      style={{
        width: "100%",
        marginTop: "1rem",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <div>
        <img src="/256.png" alt="" style={{ height: "25rem", width: "auto" }} />

        <h1>bigmoney29</h1>
        <p>
          bigmoney29 is a simple self-hosted grading server for you and your
          programming students.
        </p>
        <Collapse appear={false} in={state.errorMessage || state.fetching}>
          <div>
            {state.errorMessage ? (
              <Alert variant="danger" onClick={() => dispatch(clearError())}>
                Error logging in: {state.errorMessage}
              </Alert>
            ) : (
              <Alert variant="primary" onClick={() => dispatch(clearError())}>
                Logging in...
              </Alert>
            )}
          </div>
        </Collapse>
        <Form
          onSubmit={(e) => {
            e.preventDefault();
            dispatch(login(e.currentTarget));
          }}
        >
          <Form.Group className="mb-3" controlId="formUsername">
            <Form.Label>Username</Form.Label>
            <Form.Control
              type="text"
              name="username"
              placeholder="Josie Bruin"
            />
          </Form.Group>
          <Form.Group className="mb-3" controlId="formPassword">
            <Form.Label>Password</Form.Label>
            <Form.Control type="password" name="password" />
          </Form.Group>
          <Button variant="primary" type="submit">
            Log In
          </Button>
        </Form>
      </div>
    </div>
  );
};

export default Home;
