import { Slide, Image} from 'spectacle';
import diagram from '../assets/ingress_controller.svg';


export default function Page() {
  return (
      <Slide backgroundColor="#fff">
      <Image src={diagram} height="100%" />
    </Slide>
  );
}
