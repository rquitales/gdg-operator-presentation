import { Slide, Heading, Text } from 'spectacle';

export default function Page() {
  return (
    <Slide backgroundColor="#1F2022">
      <Heading color="#fff">"Operators"</Heading>
      <Text color="#ebebeb" paddingTop="250px" textAlign="right">
        "An operator is a K8s controller - specific to operating an
        application."
      </Text>

      <Text color="#ffda75" textAlign="center">
        Kubernetes-native application
      </Text>
    </Slide>
  );
}
