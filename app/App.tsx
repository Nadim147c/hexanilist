import { Button } from '@/components/ui/button';
import { Card, CardContent, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';

export default function App() {
  return (
    <main className="flex h-screen w-screen items-center justify-center">
      <Card className="h-max p-3">
        <CardTitle>Welcome to Hexanilist!</CardTitle>
        <CardContent>
          <Input />
          <Button>Generate</Button>
        </CardContent>
      </Card>
    </main>
  );
}
