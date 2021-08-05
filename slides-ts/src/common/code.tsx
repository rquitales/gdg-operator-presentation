import { Component, Fragment } from 'react';
import AceEditor from 'react-ace';
import styles from './code.module.css';
import './code.css';
import { Rnd } from 'react-rnd';
import { Payload } from './interfaces';

import {
  Button,
  ButtonGroup,
  Grid,
  FormGroup,
  FormControlLabel,
  Switch,
  Snackbar,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
  Card,
  CardContent,
  Box,
} from '@material-ui/core';
import MuiAlert from '@material-ui/lab/Alert';

import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/mode-sh';
import 'ace-builds/src-noconflict/mode-json';
import 'ace-builds/src-noconflict/mode-python';
import 'ace-builds/src-noconflict/theme-mono_industrial';
import 'ace-builds/src-noconflict/ext-language_tools';

interface Props {
  value?: string;
  client: WebSocket;
  runner?: string;
  pID?: string; // Hack to share pID across multiple websocket connections - mainly for terraform operations where state matters.
  outputPane?: boolean;
  language?: string;
  path?: string;
}

interface State {
  editorText: string;
  pID: string;
  responses?: Payload[]; // we should type the response struct too...
  kind: string;
  showOutput: boolean;
  autoScroll: boolean;
  snackBar: boolean;
  lastResponseType?: string;
  outputPane: boolean;
  language: string;
  path?: string;
}

function sanitise(str: string) {
  return str.replace(/&/g, '&amp;').replace(/</g, '&lt;');
}

class Editor extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      editorText: this.props.value || '',
      pID: this.props.pID || `${Math.floor(Math.random() * 9999999)}`,
      kind: this.props.runner || 'run',
      showOutput: this.props.outputPane || false,
      autoScroll: true,
      snackBar: false,
      outputPane: this.props.outputPane || false,
      language: this.props.language || 'golang',
      path: this.props.path || '',
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.handleAutoScroll = this.handleAutoScroll.bind(this);
    this.scrollToBottom = this.scrollToBottom.bind(this);
    this.toggleSnackbar = this.toggleSnackbar.bind(this);
    this.snackbarSuccess = this.snackbarSuccess.bind(this);
    this.showOutput = this.showOutput.bind(this);
  }

  killCode = () => {
    this.props.client.send(
      JSON.stringify({
        Id: this.state.pID,
        Kind: 'kill',
      } as Payload)
    );
  };

  runCode = () => {
    this.setState({
      responses: [] as Payload[],
      showOutput: true,
      snackBar: false,
    } as State);

    this.props.client.send(
      JSON.stringify({
        Id: this.state.pID,
        Kind: this.state.kind === 'shell' ? 'run' : this.state.kind,
        Path: this.state.path,
        Body:
          this.state.kind === 'shell'
            ? '#!/bin/sh\n' + this.state.editorText
            : this.state.editorText,
      } as Payload)
    );
  };

  onChange = (newValue: string) => {
    this.setState({ editorText: newValue } as State);
  };

  clear = () => {
    this.setState({ responses: [] as Payload[], snackBar: false } as State);
  };

  handleChange(event: React.ChangeEvent<{ value: unknown }>) {
    this.setState({ kind: event.target.value, snackBar: false } as State);
  }

  handleClose() {
    const toggle = !this.state.showOutput;
    this.setState({ showOutput: toggle } as State);
  }

  handleAutoScroll() {
    this.setState({ autoScroll: !this.state.autoScroll } as State);
  }

  seqID = 0;

  componentDidUpdate() {
    this.scrollToBottom();
  }

  scrollToBottom() {
    if (!this.state.autoScroll) {
      return;
    }
    this.output?.scrollIntoView({ behavior: 'smooth' });
  }

  output = null as HTMLDivElement | null;

  showOutput() {
    return (
      <Fragment>
        <pre>
          {(this.state.responses ?? [])
            .filter((resp) => {
              return resp.Body ?? false;
            })
            .map((resp) => {
              return (
                <div
                  key={resp.Id + resp.Kind + `${resp.seqID}`}
                  className={resp.Kind}
                  ref={(ele) => {
                    this.output = ele;
                  }}
                >
                  {sanitise(resp.Body!)}
                </div>
              );
            })}
        </pre>
        <FormGroup row>
          <FormControlLabel
            control={
              <Switch
                checked={this.state.autoScroll}
                onChange={this.handleAutoScroll}
                name="enableAutoScroll"
                color="primary"
                size="small"
              />
            }
            label={<Box style={{ color: '#fff' }}>Auto-scroll</Box>}
          />
          <Button
            variant="outlined"
            color="secondary"
            size="small"
            onClick={this.clear}
          >
            Clear
          </Button>
        </FormGroup>
      </Fragment>
    );
  }

  toggleSnackbar(event?: React.SyntheticEvent, reason?: string) {
    if (reason === 'clickaway') {
      return;
    }

    this.setState({ snackBar: false } as State);
  }

  snackbarSuccess() {
    const lastResp = this.state.responses![this.state.responses!.length - 1];
    if (!lastResp) {
      return false;
    }

    return lastResp.Body === '';
  }

  showSnackBar() {
    return (
      <div>
        <Snackbar
          open={this.state.snackBar && this.snackbarSuccess()}
          onClose={this.toggleSnackbar}
        >
          <MuiAlert
            elevation={6}
            variant="filled"
            onClose={this.toggleSnackbar}
            severity="success"
          >
            {`${this.state.kind}`} finished!
          </MuiAlert>
        </Snackbar>
        {/* {/* <Alert severity="error">This is an error message!</Alert> */}
        <Snackbar
          open={this.state.snackBar && !this.snackbarSuccess()}
          onClose={this.toggleSnackbar}
        >
          <MuiAlert
            elevation={6}
            variant="filled"
            onClose={this.toggleSnackbar}
            severity="error"
          >
            Error running "{`${this.state.kind}`}" command!
          </MuiAlert>
        </Snackbar>
      </div>
    );
  }

  showEditor() {
    return (
      <Grid container xs={12} spacing={1}>
        <Grid item xs={12}>
          <AceEditor
            className={styles.ace_editor}
            mode={this.state.language}
            theme="mono_industrial"
            name="example"
            onChange={this.onChange}
            style={{
              width: `100%`,
              height: '400px',
              padding: `0`,
              margin: `0`,
            }}
            setOptions={{
              showGutter: true,
              enableBasicAutocompletion: false,
              enableSnippets: false,
              enableLiveAutocompletion: false,
              showLineNumbers: true,
              tabSize: 2,
              hasCssTransforms: true,
            }}
            value={this.state.editorText}
            editorProps={{ $blockScrolling: true }}
          />
        </Grid>
      </Grid>
    );
  }

  showControls() {
    console.log(this.state);
    return (
      <Grid container spacing={6}>
        {this.state.kind !== 'saveFile' ? (
          <Grid item xs={4}>
            <ButtonGroup
              variant="contained"
              size="large"
              color="primary"
              style={{ display: 'block !important' }}
            >
              <Button onClick={this.runCode} color="primary">
                Run
              </Button>
              <Button
                onClick={() => {
                  this.setState({ editorText: this.props.value || '' });
                }}
                color="primary"
              >
                Reset Code
              </Button>
              <Button onClick={this.killCode} color="secondary">
                Kill
              </Button>
              {!this.state.outputPane ? (
                <Button onClick={this.handleClose}>Output</Button>
              ) : null}
            </ButtonGroup>
          </Grid>
        ) : (
          <Grid item xs={3}>
            <Button onClick={this.runCode} color="primary" variant="contained">
              Save
            </Button>
            <Button
                onClick={() => {
                  this.setState({ editorText: this.props.value || '' });
                }}
                color="primary"
              >
                Reset file
              </Button>
          </Grid>
        )}

        {this.state.kind !== 'saveFile' ? (
          <Grid item xs={2}>
            <FormControl variant="filled" size="small">
              <InputLabel id="select-runner-label">Runner</InputLabel>
              <Select
                labelId="select-runner-label"
                id="select-runner"
                value={this.state.kind}
                onChange={this.handleChange}
                variant="filled"
              >
                <MenuItem value={'run'}>Golang</MenuItem>
                <MenuItem value={'shell'}>Shell</MenuItem>
                <MenuItem value={'kubectlApply'}>kubectl - Apply</MenuItem>
                <MenuItem value={'kubectlCreate'}>kubectl - Create</MenuItem>
                <MenuItem value={'kubectlDelete'}>kubectl - Delete</MenuItem>
                <MenuItem value={'terraformApply'}>terraform - apply</MenuItem>
                <MenuItem value={'terraformDestroy'}>
                  terraform - destroy
                </MenuItem>
              </Select>
            </FormControl>
          </Grid>
        ) : null}
      </Grid>
    );
  }

  render() {
    this.props.client.onmessage = (message) => {
      const dataFromServer: Payload = JSON.parse(message.data);
      dataFromServer.seqID = this.seqID;
      this.seqID += 1;
      let stateToChange = {
        responses: [...(this.state.responses ?? []), dataFromServer],
      };

      this.setState({
        ...stateToChange,
        snackBar: dataFromServer.Kind === 'end',
        lastResponseType: dataFromServer.Kind,
      } as State);

      console.log(this.state);
    };

    return (
      <Card>
        <CardContent>
          {this.state.outputPane ? (
            <Grid container spacing={2}>
              <Grid item xs={6}>
                {this.showEditor()}
              </Grid>
              <Grid item xs={6}>
                <div className="outputFixed">{this.showOutput()}</div>
              </Grid>

              {this.showControls()}

              {this.showSnackBar()}
            </Grid>
          ) : (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                {this.showEditor()}
              </Grid>
              {this.showControls()}

              {this.state.kind !== 'saveFile' ? (
                <Rnd
                  style={{ display: this.state.showOutput ? 'block' : 'none' }}
                  default={{
                    x: 0,
                    y: 0,
                    width: 580,
                    height: 520,
                  }}
                  className="output"
                >
                  {this.showOutput()}
                </Rnd>
              ) : null}
              {this.showSnackBar()}
            </Grid>
          )}
        </CardContent>
      </Card>
    );
  }
}

export default Editor;
