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
      <Text color="black">
        Do:
        <UnorderedList color="#black">
          <ListItem>Application has declarative API</ListItem>
          <ListItem>Scope down to namespace/cluster level</ListItem>
          <ListItem>Top-level kubectl/API support</ListItem>
          <ListItem>Encapsulates business logic</ListItem>
        </UnorderedList>
        Don't:
        <UnorderedList color="#black">
          <ListItem>Application has declarative API</ListItem>
          <ListItem>Scope down to namespace/cluster level</ListItem>
          <ListItem>Top-level kubectl/API support</ListItem>
          <ListItem>Encapsulates business logic</ListItem>
        </UnorderedList>
      </Text>
    </Slide>
  );
}
