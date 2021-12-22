/**
 * Data-rich summary of your scores and stuff.
 */

import Card from 'react-bootstrap/Card';
import Stack from 'react-bootstrap/Stack';

const Summary = (props) => {
  return (
    <Stack direction="horizontal" >
      <Card style={{ width: "18rem" }}>
        <Card.Title>Assignment</Card.Title>
        <Card.Body>
          Highest Score
        </Card.Body>
      </Card>
    </Stack>
  );
};

export default Summary;
