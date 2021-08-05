import { Slide, Heading, FlexBox, Box, Image, Text } from 'spectacle';
import logo from '../assets/title.gif';

export default function Page() {
  return (
    <Slide backgroundColor="#1F2022">
      <FlexBox justifyContent="center" position="absolute" width={0.92}>
        <Box>
          <Heading color="#fff">Intro to Kubernetes Operators</Heading>
          <Text textAlign="center" color="grey">Ramon Quitales</Text>
        </Box>
        <Box>
        <Image src={logo} />
        </Box>
      </FlexBox>
    </Slide>
  );
}
