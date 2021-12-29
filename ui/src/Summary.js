/**
 * Data-rich summary of your scores and stuff.
 */

import Card from "react-bootstrap/Card";
import Container from "react-bootstrap/esm/Container";
import Stack from "react-bootstrap/Stack";
import Button from "react-bootstrap/Button";
import Alert from "react-bootstrap/Alert";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";

const Summary = () => {
  const navigate = useNavigate();
  const { 0: cookies } = useCookies(["jwt"]);

  if (!cookies.jwt)
    return <Alert variant="danger">You are not logged in.</Alert>;

  return (
    <Container>
      <Button variant="link" onClick={() => navigate("/me")}>
        ⮐ Back to assignments
      </Button>
      <Stack direction="horizontal">
        <Card style={{ width: "18rem" }}>
          <Card.Title>Assignment</Card.Title>
          <Card.Body>Highest Score</Card.Body>
        </Card>
      </Stack>
    </Container>
  );
};

export default Summary;
