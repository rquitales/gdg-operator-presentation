import { Grid } from '@material-ui/core';
import { Slide, Heading } from 'spectacle';
import Term from '../common/terminal';
import raw from 'raw.macro';
import Editor from '../common/code';

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);
// const client = new WebSocket('ws://localhost:8082/socket');

const startCode = raw('./code/deploy_run.txt');
const path = '/home/nonroot/projects/serverless/api/v1alpha1/function_types.go';

export default function Page() {
  return (
    <Slide>
      <Heading>
        Run Operator
      </Heading>
      <Grid container xs={12} spacing={1}>
        <Grid item xs={6}>
          <Editor
            client={client}
            value={startCode}
            language={'sh'}
            runner="shell"
            outputPane={false}
            path={path}
          />
        </Grid>
        <Grid item xs={6}>
          <Term />
        </Grid>
      </Grid>
    </Slide>
  );
}
