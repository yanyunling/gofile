interface Log {
  id: string;
  title: string;
  content: string;
  level: string;
  username: string;
  createTime: number;
}

interface LogCondition {
  title: string;
  content: string;
  level: string;
  username: string;
  startTime: number;
  endTime: number;
}
