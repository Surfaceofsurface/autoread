<script lang="ts">
  import type { User } from "./types/users";

  export let user: User;

  const time = new Date();
  $: {
    console.log(user);
  }
</script>

<div
  class="
    bg-slate-300 
    rounded-md 
    p-8
    h-full
    md:p-4
    "
>
  <header class="grid rounded grid-cols-3 gap-4 md:gap-2 h-full">
    <ul class="flex flex-col justify-around h-full">
      <li class="pb-3">
        <div class="text-lg md:text-base font-semibold">账号</div>
        <span class="font-mono">{user.Username}</span>
      </li>
      <li class="pb-3">
        <div class="text-lg md:text-base font-semibold">开始时间</div>
        <span class="font-mono">{time}</span>
      </li>
      <li class="pb-3">
        <div class="text-lg md:text-base font-semibold">并发数</div>
        <span class="font-mono">2</span>
      </li>
      <li class="pb-3">
        <div class="text-lg md:text-base font-semibold">日志</div>
        {#if user.Logs && user.Logs.length > 0}
          {#each user.Logs as log}
            <div class="font-mono">{log}</div>
          {/each}
        {:else}
          <p class="font-mono">暂无日志</p>
        {/if}
      </li>
    </ul>
    <article class="flex flex-col  h-full col-span-2">
      <ul>
        {#each user.PendingBook as bk}
          <div class="py-1">
            <li class="text-gray-500 text-sm flex justify-between ">
              <span>{bk.Title}</span>
              <i>{bk.Process ? (100 * bk.Process).toFixed(2) : "??"}%</i>
            </li>
            <li>
              <div class="w-full h-1.5 ">
                <div
                  class="
                  transition-all
                    h-full
                    w
              bg-gradient-to-r
             from-blue-300
             to-indigo-500
            "
                  style="width:{bk.Process ? 100 * bk.Process : 1}%;"
                />
              </div>
            </li>
          </div>
        {/each}
      </ul>
    </article>
  </header>
</div>
