import {
  Slide,
  Heading,
  Box,
  Image,
  Text,
  UnorderedList,
  ListItem,
} from 'spectacle';
import logo from '../assets/k8s.png';

export default function Page() {
  return (
    <Slide backgroundColor="#1F2022">
      <Heading color="#fff">K8s in a minute...</Heading>
      <Text color="#fff">
        <UnorderedList color="#fff">
          <ListItem>
            Container orchestration
          </ListItem>
          <ListItem>
            Scale workloads across clusters
          </ListItem>
          <ListItem>
            API centric
          </ListItem>
          <ListItem>
            Native stable resources (Deployments, Pods etc)
          </ListItem>
        </UnorderedList>
      </Text>
      <Box paddingLeft="1000px">
        <Image src={logo} width="50%" />
      </Box>
    </Slide>
  );
}
