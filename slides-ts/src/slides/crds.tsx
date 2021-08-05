import { Slide, Heading, Text, UnorderedList, ListItem } from 'spectacle';

export default function Page() {
  return (
    <Slide backgroundColor="#1F2022">
      <Heading color="#fff">Extending K8s - CRDs</Heading>
      <Text color="#fff">
        CustomResourceDefinitions allow users to create new types of resources
        without adding another API server.
        <UnorderedList color="#fff">
          <ListItem>Easy to create/define</ListItem>
          <ListItem>Top level support with kubectl</ListItem>
        </UnorderedList>
      </Text>
    </Slide>
  );
}
