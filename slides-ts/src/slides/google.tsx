import { Slide, FlexBox, Box, Image, Text, UnorderedList, ListItem } from 'spectacle';
// import logo from '../assets/google-logo.png';
import logogcp from '../assets/gcp.png';

export default function Page() {
  return (
    <Slide backgroundColor="#fff" textColor="#1F2022">
      <FlexBox paddingTop="100px">
        <Box>
          <Image src={logogcp} size="70%"/>
        </Box>
        <Box>
          <Text>
              <UnorderedList>
    <ListItem>
        Engineer @ Google
    </ListItem>
    <ListItem>
        Lead development of Google Cloud Anthos products
    </ListItem>
    <ListItem>
        Golang lover
    </ListItem>
              </UnorderedList>
          </Text>
        </Box>
      </FlexBox>
      <Text textAlign="center" textColor="#1F2022" fontSize="16px">
        *The contents of this slide/presentation are my own personal views and
        beliefs.
      </Text>
    </Slide>
  );
}
