export type Allusers = {
  [uid: string]: User;
};
export type User = {
  Username: string;
  PendingBook: {
    Title: string;
    Process?: number;
  }[];
  Logs?: string[];
};
export type WsProcessDTO = {
  who: string;
  d: number;
  type: string;
  bid: string;
};
export type WsDoneDTO = {
  who: string;
  type: string;
  bid: string;
};
