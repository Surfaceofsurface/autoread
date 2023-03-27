<script lang="ts">
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
  type user = {
    psw?: string;
    uid?: string;
    conc: number;
    gap: number;
  };

  const userlist: user = { conc: 2, gap: 3 };
  let canUpload: boolean = true;
  async function userUpload() {
    canUpload = false;
    await fetch("/api/user", {
      method: "post",
      body: JSON.stringify(userlist),
    });
    dispatch("userUpdate");
    canUpload = true;
    userlist.psw = "";
    userlist.uid = "";
    userlist.conc = 2;
    userlist.gap = 3;
  }
</script>

<div
  class="
    flex
    flex-col
    justify-around
    rounded-md 
    p-8
    md:p-4
    outline-dashed outline-2 outline-offset-2 outline-blue-500
    "
>
  <div class=" text-xl text-blue-500 ">新建账号</div>
  <div>
    <div class="col-span-6 sm:col-span-3">
      <label
        for="账号"
        class="block text-sm font-medium leading-6 text-gray-900">账号</label
      >
      <input
        bind:value={userlist.uid}
        type="text"
        class="font-mono mt-2 block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
      />
    </div>

    <div class="col-span-6 sm:col-span-3">
      <label
        for="last-name"
        class=" block text-sm font-medium leading-6 text-gray-900">密码</label
      >
      <input
        bind:value={userlist.psw}
        type="password"
        name="last-name"
        id="last-name"
        class="font-mono mt-2 block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
      />
    </div>
    <div class="col-span-6 sm:col-span-3">
      <label
        for="last-name"
        class=" block text-sm font-medium leading-6 text-gray-900">并发数</label
      >
      <input
        bind:value={userlist.conc}
        type="number"
        name="last-name"
        id="conc"
        class="font-mono mt-2 block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
      />
    </div>
    <div class="col-span-6 sm:col-span-3">
      <label
        for="last-name"
        class=" block text-sm font-medium leading-6 text-gray-900"
        >平均翻页时间</label
      >
      <input
        bind:value={userlist.gap}
        type="number"
        name="last-name"
        id="gap"
        class="font-mono mt-2 block w-full rounded-md border-0 p-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
      />
    </div>
  </div>
  <div class=" px-4 py-3 text-right sm:px-6">
    <button
      disabled={!canUpload}
      on:click={userUpload}
      type="submit"
      class="inline-flex justify-center rounded-md bg-indigo-600 py-2 px-3 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"
      >提交
      <span />
      <i class="{canUpload ? 'hidden' : 'inline'} ">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          style="width:24px;height:24px"
          viewBox="0 0 20 20"
        >
          <circle
            id="spinner"
            cx="10"
            cy="10"
            r="6"
            fill="none"
            stroke="white"
            stroke-width="2"
            stroke-dasharray="50.2 50"
          />
          <animate
            attributeName="stroke-dashoffset"
            values="52;0;52"
            dur="5s"
            repeatCount="indefinite"
          />
          <animateTransform
            attributeName="transform"
            type="rotate"
            values="0 0 0;360 0 0"
            dur="9s"
            repeatCount="indefinite"
          />
        </svg>
      </i>
    </button>

    <div />
  </div>
</div>
