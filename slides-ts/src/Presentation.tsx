import { Deck, FlexBox, Box, FullScreen, Progress } from 'spectacle';
import TitleS from './slides/01-title';
// import Slide2 from './slides/02-about';
import CreateK8s from './slides/03-tf';
import OperatorChoices from './slides/04-operator_choices';
import Scaffold from './slides/scaffold_kubebuilder';
import Slide6 from './slides/test_cli';
import Slide7 from './slides/create_api';
import Slide8 from './slides/create_controller';
import Google from './slides/google';
import ButFirstK8s from './slides/but_first';
import K8sIntro from './slides/K8s-intro';
import K8sDiagram from './slides/k8s-diagram';
import Controllers from './slides/controllers';
import CRDs from './slides/crds';
import Operators from './slides/operators';
import WhenToUse from './slides/when_to_use';
import DeployFunc from './slides/deploy_func';
import SimpleArchitecture from './slides/serverless_architecture';
import Improvement from './slides/improvement';
import Alternatives from './slides/alternatives';
// import Teardown from './slides/teardown_cluster';

const template = () => (
  <FlexBox
    justifyContent="space-between"
    position="absolute"
    bottom={0}
    width={1}
  >
    <Box padding="0 1em">
      <FullScreen color="black" size={13} />
    </Box>
    <Box padding="1em">
      <Progress color="black" size={13} />
    </Box>
  </FlexBox>
);

const theme = {
  colors: {
    primary: '#444',
    secondary: '#1F2022',
    tertiary: '#CECECE',
    quaternary: '#CECECE',
  },
  fontSizes: {
    header: '64px',
    paragraph: '22px',
  },
  fonts: {
    header: '"Ubuntu", Helvetica, Arial, sans-serif',
    text: '"Fira sans", Helvetica, Arial, sans-serif',
  },
};

// Hack for terraform.
const randID = `${Math.floor(Math.random() * 9999999)}`;

export default function Presentation() {
  return (
    <Deck theme={theme} template={template}>
      <TitleS />
      <Google />
      <ButFirstK8s />
      {/* <Slide2 /> */}
      <CreateK8s pID={randID} />
      <K8sIntro />
      <K8sDiagram />
      <Controllers />
      <CRDs />
      <Operators />
      <WhenToUse />
      <SimpleArchitecture />
      <OperatorChoices />
      <Scaffold />
      <Slide7 />
      <Slide8 />
      <Slide6 />
      <DeployFunc />
      <Improvement />
      <Alternatives />
      {/* <Teardown pID={randID} /> */}
    </Deck>
  );
}
