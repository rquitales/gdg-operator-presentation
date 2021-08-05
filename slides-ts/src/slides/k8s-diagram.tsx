import { Slide, Image } from 'spectacle';
import diagram from '../assets/k8s-components.svg';

export default function Page() {
  return (
    <Slide backgroundColor="#fff">
      <Image src={diagram} width="100%" />
    </Slide>
  );
}
