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
      <Heading color="black">Should we use this in production?</Heading>
      <Text color="black" fontSize={"29px"}>
        <b>Maybe not...</b>
        <UnorderedList color="#black" fontSize={"24px"}>
          <ListItem>Use a managed service if low throughput (GCP Cloud Functions or Cloud Run)</ListItem>
          <ListItem>Knative (open sourced GCP Cloud Run)</ListItem>
          <ListItem>OpenFaas</ListItem>
          <ListItem>Kubeless (similar implementation)</ListItem>
        </UnorderedList>
      </Text>
    </Slide>
  );
}
