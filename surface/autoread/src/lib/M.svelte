<script lang="ts">
  import AppendUser from "./AppendUser.svelte";
  import UserCard from "./UserCard.svelte";
  import type { Allusers, WsProcessDTO } from "./types/users";
  let users: Allusers = {};

  const ws = new WebSocket("ws://" + window.location.host + "/api/per");
  ws.addEventListener("message", function (event) {
    const dt = JSON.parse(event.data);
    if (dt.type) {
      if (dt.type === "process") {
        const dto = dt as WsProcessDTO;
        users[dto.who].PendingBook.find((v) => {
          return v.Title === dto.bid;
        }).Process = dto.d;
        users[dto.who].PendingBook = users[dto.who].PendingBook;
      }
      if (dt.type === "done") {
        const dto = dt as WsProcessDTO;
        users[dto.who].Logs = users[dto.who].Logs || [];
        users[dto.who].Logs = [
          ...users[dto.who].Logs,
          dto.bid + "进度已到100%,若未显示已完成,请手动将进度置为0,重来一遍",
        ];
      }
    }
  });
  async function getAllUsers() {
    const data = await fetch("/api/alluser");
    users = await data.json();
  }
  getAllUsers();
</script>

<main
  class="
 bg-slate-200
  p-5
  min-h-screen
  grid grid-cols-1 gap-4
  md:grid-cols-2 "
>
  {#each Object.entries(users) as user}
    <UserCard user={user[1]} />
  {/each}

  <AppendUser on:userUpdate={getAllUsers} />
</main>
