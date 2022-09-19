import {
  BrowserRouter as Router,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";

import { useDispatch, useSelector } from "react-redux";

import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import Button from "react-bootstrap/Button";

import Landing from "./Landing";
import Results from "./Results";
import Assignment from "./Assignment";
import Home from "./Home";

// Navbar for the app.
const Moneybar = ({}) => {
  const dispatch = useDispatch();

  return (
    <Navbar bg="light">
      <Container>
        <Navbar.Brand>
          <img
            alt=""
            src="/256.png"
            width="30"
            height="30"
            className="d-inline-block align-top"
          />{" "}
          bigmoney29
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse>
          <Nav className="me-auto">
            <Nav.Link href="#home">Assignments</Nav.Link>
          </Nav>
        </Navbar.Collapse>
        <Navbar.Collapse className="justify-content-end">
          <Navbar.Text>
            <Button variant="danger" onClick={() => dispatch()}>
              Logout
            </Button>
          </Navbar.Text>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};

const App = () => {
  const auth = useSelector((state) => state.auth);

  return (
    <>
      {auth.token && <Moneybar loggedIn={auth.token} />}
      <Router>
        <Routes>
          <Route path="/">
            <Route index element={!auth.token ? <Home /> : <Landing />} />
            <Route path="assignment/:id" element={<Assignment />} />
            <Route path="results/:id" element={<Results />} />
          </Route>
        </Routes>
      </Router>
    </>
  );
};

export default App;
