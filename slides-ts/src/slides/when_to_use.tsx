import {
  Slide,
  Heading,
  Text,
  UnorderedList,
  ListItem,
} from 'spectacle';

export default function Page() {
  return (
    <Slide backgroundColor="#fff0b8">
      <Heading color="black">When to Use?</Heading>
      <Text color="black" fontSize={"29px"}>
        <b>Do:</b>
        <UnorderedList color="#black" fontSize={"24px"}>
          <ListItem>Application has declarative API</ListItem>
          <ListItem>Scope down to namespace/cluster level</ListItem>
          <ListItem>Top-level kubectl/API support</ListItem>
          <ListItem>Encapsulates business logic</ListItem>
        </UnorderedList>
        <hr />
        <b>Don't:</b>
        <UnorderedList color="#black" fontSize={"24px"}>
          <ListItem>Config file is used for application configuration</ListItem>
          <ListItem>No need to reconcile state</ListItem>
        </UnorderedList>
        <Text textAlign="center" fontSize="21px" color="primary"> Use ConfigMap or Secrets instead</Text>
      </Text>
    </Slide>
  );
}
