/**
 * Check when your stuff is due, or edit your class.
 */

import { useState } from 'react';

import JWTContext from './jwt-context';

import Container from 'react-bootstrap/Container';
import Stack from 'react-bootstrap/Stack';
import ProgressBar from 'react-bootstrap/ProgressBar';
import Spinner from 'react-bootstrap/Spinner';

const Me = () =>
  <JWTContext.Consumer>
    {({ jwt, assignments, updateAssignments }) => {
      if (!assignments)
        updateAssignments();
      return (
        <Container>
          {!assignments ?
            <Spinner animation="border" role="status">
              <span className='visually-hidden'>Loading...</span>
            </Spinner> :

            <Stack direction="vertical">
              {[0, 1, 1, 1].map((_, idx) => (
                <Container key={idx}>
                  <h2>Name</h2>
                  <ProgressBar now={++idx * 10} />
                </Container>
              ))}
            </Stack>
          }
        </Container>
      );
    }
    }
  </JWTContext.Consumer>

export default Me;
