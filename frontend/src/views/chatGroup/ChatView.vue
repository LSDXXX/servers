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
          <div v-html="content(message.content)"></div>
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
    this.socket = new WebSocket("ws://" + window.location.host + "/api/chat");
    this.socket.onopen = function () {
      console.log("connected to server");
    };
    var chat = this;
    this.socket.onmessage = function (event) {
      console.log(event);
      event.data.text().then(function (data) {
        var obj = JSON.parse(data);
        var found = false;
        for (var i = 0; i < chat.messages.length; i++) {
          if (chat.messages[i].id == obj.message.id) {
            chat.messages[i].content = obj.message.content.parts[0];
            found = true;
            break;
          }
        }
        if (!found) {
          chat.messages.push({
            id: obj.message.id,
            content: obj.message.content.parts[0],
            from: "bot",
          });
        }
      });
      // this.messages.push(data)
    };
  },
  methods: {
    send() {
      if (this.message) {
        // 发送消息到服务端
        this.socket.send(this.message);
        this.messages.push({
          id: "default",
          from: "me",
          content: this.message,
        });
        this.message = "";
      }
    },
    content(data) {
      const marked = require("marked");
      return marked.marked(data);
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
