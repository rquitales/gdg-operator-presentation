import { Slide, Image} from 'spectacle';
import diagram from '../assets/simple_architecture.svg';


export default function Page() {
  return (
      <Slide backgroundColor="#fff">
      <Image src={diagram} height="100%" />
    </Slide>
  );
}
