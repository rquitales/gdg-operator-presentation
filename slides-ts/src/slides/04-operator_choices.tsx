import {Slide, Heading, OrderedList, ListItem, Text, UnorderedList} from 'spectacle';


export default function Page() {
  return (
    <Slide>
      <Heading>
        Let's Build!
      </Heading>
      <Text fontSize="22px">
        It's an application: build using your <b>any</b> preferred language and/or framework
      </Text>
      <OrderedList fontSize="28px">
        <ListItem>
          kubebuilder - K8s SIG
        </ListItem>
        <ListItem>
          Operator SDK - Operator Framework
          <UnorderedList fontSize="26px">
            <ListItem>
              Ansible
            </ListItem>
            <ListItem>
              Helm
            </ListItem>
            <ListItem>
              Go
            </ListItem>
          </UnorderedList>
        </ListItem>
        <ListItem>
          Operator Framework: KOPF (Python)
        </ListItem>
        <ListItem>
          Roll your own
        </ListItem>
      </OrderedList>
    </Slide>
  );
}
