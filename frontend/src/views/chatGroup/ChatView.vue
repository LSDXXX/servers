<template>
  <div>
    <div class="chat">
      <div
        v-for="(message, index) in messages"
        :key="index"
        :class="{ me: message.from === 'me', other: message.from !== 'me' }"
      >
        <div class="message">
          <div class="from">{{ message.from }}</div>
          <div class="content">{{ message.content }}</div>
        </div>
      </div>
    </div>
    <div class="input">
      <input type="text" v-model="message" @keyup.enter="send" />
      <button @click="send">发送</button>
    </div>
  </div>
</template>

<script>
import io from "socket.io-client";

export default {
  data() {
    return {
      messages: [],
      message: "",
      socket: null,
    };
  },
  created() {
    // 建立 WebSocket 连接
    let url = "ws://" + window.location.host;
    console.log(url);
    this.socket = io(url, { path: "/api/chat" });
    this.socket.on("connect", () => {
      console.log("connected to server");
    });
    // 监听服务端推送的消息
    this.socket.on("message", (data) => {
      this.messages.push(data);
    });
  },
  methods: {
    send() {
      if (this.message) {
        // 发送消息到服务端
        this.socket.emit("message", {
          from: "me",
          content: this.message,
        });
        this.messages.push({
          from: "me",
          content: this.message,
        });
        this.message = "";
      }
    },
  },
};
</script>

<style>
.chat {
  display: flex;
  flex-direction: column;
}

.message {
  margin: 5px;
  padding: 10px;
  border-radius: 5px;
  background-color: #eee;
}

.me .from {
  text-align: right;
}

.other .from {
  text-align: left;
}

.input {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

input[type="text"] {
  width: 80%;
  margin-right: 10px;
  padding: 5px;
}

button {
  padding: 5px 10px;
  border-radius: 5px;
  background-color: #4caf50;
  color: #fff;
  border: none;
  cursor: pointer;
}
</style>
