import { Slide, Heading } from 'spectacle';
import Editor from '../common/code';
import { Grid } from '@material-ui/core';

let startCode = `#!/bin/bash
cd /Users/rquitales/code/faas
make run`;

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);

export default function Page(props: any) {
  return (
    <Slide backgroundColor="white">
      <Heading>Start</Heading>
      <Grid item xs={12} spacing={1}>
        <Editor value={startCode} client={client} />
      </Grid>
    </Slide>
  );
}
