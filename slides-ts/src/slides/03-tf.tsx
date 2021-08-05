import { Slide, Heading, Text } from 'spectacle';
import Editor from '../common/code';

let startCode = `# INSERT YOUR TOKEN BELOW

terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "2.4.0"
    }
  }
}

variable "do_token" {
    default = <INSERT_TOKEN_HERE>
}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_kubernetes_versions" "v1_20" {
  version_prefix = "1.20."
}

resource "digitalocean_kubernetes_cluster" "personal-k8s" {
  name   = "k8s-operator-example"
  region = "sfo3"
  version = data.digitalocean_kubernetes_versions.v1_20.latest_version

  node_pool {
    name       = "worker-pool-1"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}

resource "digitalocean_project" "k8s-operator-example" {
  name        = "Kubernetes Operator Example"
  description = "A project to demo a simple K8s operator."
  environment = "Development"
}

resource "digitalocean_project_resources" "k8s-resources" {
  project = digitalocean_project.k8s-operator-example.id
  resources = ["do:kubernetes:\${digitalocean_kubernetes_cluster.personal-k8s.id}"]
}
`;

const host = window.location.host;
const client = new WebSocket(`ws://${host}/socket`);

export default function Page(props: any) {
  return (
    <Slide>
      <Heading>Create K8s Cluster</Heading>
      <Editor
        value={startCode}
        client={client}
        runner="terraformApply"
        pID={props.pID}
        language="json"
        outputPane={true}
      />
      <Text fontSize="12px">
        Rember to switch to terraformDestroy if you want to destroy these clusters after testing!
      </Text>
    </Slide>
  );
}
