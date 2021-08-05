import {
  Slide,
  Heading,
  Text,
  UnorderedList,
  ListItem,
} from 'spectacle';

export default function Page() {
  return (
    <Slide backgroundColor="#1F2022">
      <Heading color="#fff">What are controllers?</Heading>
      <Text color="#fff">
        <UnorderedList color="#fff">
          <ListItem>Infinite loop to course correct</ListItem>
          <UnorderedList color="grey">
            <ListItem>Regulates state of system</ListItem>
            <ListItem>Similar to heat pumps</ListItem>
          </UnorderedList>
          <ListItem>Reconciles given state</ListItem>
          <UnorderedList color="grey">
            <ListItem>Watches K8s cluster state via API server</ListItem>
            <ListItem>Example: ReplicaSet controller</ListItem>
          </UnorderedList>
        </UnorderedList>
      </Text>
    </Slide>
  );
}
