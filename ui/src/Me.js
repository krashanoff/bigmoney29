/**
 * Check when your stuff is due, or edit your class.
 */

import { useEffect, useState } from 'react';
import { useCookies } from 'react-cookie';

import Alert from 'react-bootstrap/Alert';
import Placeholder from 'react-bootstrap/Placeholder';
import Container from 'react-bootstrap/Container';
import Stack from 'react-bootstrap/Stack';
import ProgressBar from 'react-bootstrap/ProgressBar';
import Spinner from 'react-bootstrap/Spinner';
import Button from 'react-bootstrap/Button';
import { useNavigate } from 'react-router-dom';

const Me = () => {
  const navigate = useNavigate();
  const { 0: cookies } = useCookies(["jwt"]);
  const [errors, setErrors] = useState({ auth: "" });
  const [assignments, setAssignments] = useState([]);

  // Get user overview.
  useEffect(() => {
    fetch("http://localhost:8081/cash/assignments", {
      method: "get",
      headers: {
        "Authorization": `Bearer: ${cookies.jwt}`,
      },
    })
      .then(r => {
        if (r.status < 300 && r.status >= 200)
          return { data: r.json() };
        return { status: r.status };
      })
      .then(obj => {
        if (obj.status) {
          console.info("Status is", obj.status);
          return setErrors({ auth: "You are not logged in." });
        }
        setAssignments(obj.data);
      })
      .catch(errorMessage => setErrors({ auth: errorMessage }));
    return () => setAssignments([]);
  }, [cookies]);

  // Loading page.
  if (!assignments)
    return (
      <Spinner animation="border" variant="primary" role="status">
        <span className="visually-hidden">Loading...</span>
      </Spinner>
    );

  return (
    <Container>
      {errors.auth && <Alert variant="danger">{errors.auth}</Alert>}
      {!errors.auth && !assignments ?
        <h2 aria-hidden="true">
          <Placeholder xs={6} />
        </h2> :

        <Stack direction="vertical">
          {[0, 1, 1, 1].map((assignmentId, idx) => (
            <Container key={idx}>
              <Stack direction="horizontal">
                <h2>Assignment {idx + 1}</h2>
                <Button className='ms-auto' variant="primary" onClick={() => {
                  navigate(`/summary/${assignmentId}`);
                }}>View</Button>
              </Stack>
              <ProgressBar now={++idx * 10} />
            </Container>
          ))}
        </Stack>
      }
    </Container>
  );
}

export default Me;
