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
          <div class="from" v-html="content(message.content)"></div>
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
import Clipboard from "clipboard";
export default {
  mounted() {
    this.$nextTick(() => {
      this.clipboard = new Clipboard(".copy-btn");
      // 复制成功失败的提示
      this.clipboard.on("success", () => {
        this.$message.success("复制成功");
      });
      this.clipboard.on("error", () => {
        this.$message.error("复制成功失败");
      });
    });
  },
  data() {
    return {
      messages: [],
      message: "",
      socket: null,
      clipboard: null,
    };
  },
  created() {
    // 建立 WebSocket 连接
    this.socket = new WebSocket("ws://" + window.location.host + "/api/chat");
    this.socket.onopen = function () {
      console.log("connected to server");
    };
    var chat = this;
    this.socket.onerror = function () {
      alert("ws connect error");
      chat.$router.push("/");
    };
    this.socket.onclose = function () {
      location.reload();
    };
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
      var MarkdownIt = require("markdown-it");
      var hljs = require("highlight.js");
      var md = new MarkdownIt({
        html: true,
        linkify: true,
        breaks: true,
        typographer: true,
        highlight: function (str, lang) {
          // 当前时间加随机数生成唯一的id标识
          const codeIndex =
            parseInt(Date.now()) + Math.floor(Math.random() * 10000000);
          // 复制功能主要使用的是 clipboard.js
          let html = `<button class="copy-btn" type="button" data-clipboard-action="copy" data-clipboard-target="#copy${codeIndex}">复制</button>`;
          const linesLength = str.split(/\n/).length - 1;
          // 生成行号
          let linesNum = '<span aria-hidden="true" class="line-numbers-rows">';
          for (let index = 0; index < linesLength; index++) {
            linesNum = linesNum + "<span></span>";
          }
          linesNum += "</span>";

          if (lang && hljs.getLanguage(lang)) {
            try {
              // highlight.js 高亮代码
              const preCode = hljs.highlight(lang, str, true).value;
              html = html + preCode;
              if (linesLength) {
                html += '<b class="name">' + lang + "</b>";
              }
              // 将代码包裹在 textarea 中，由于防止textarea渲染出现问题，这里将 "<" 用 "&lt;" 代替，不影响复制功能
              return `<pre class="hljs"><code>${html}</code>${linesNum}</pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(
                /<\/textarea>/g,
                "&lt;/textarea>"
              )}</textarea>`;
            } catch (error) {
              console.log(error);
            }
          }
          const preCode = md.utils.escapeHtml(str);
          html = html + preCode;
          return `<pre class="hljs"><code>${html}</code>${linesNum}</pre><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy${codeIndex}">${str.replace(
            /<\/textarea>/g,
            "&lt;/textarea>"
          )}</textarea>`;
        },
      });
      return md.render(data);
      // const marked = require("marked");
      // return marked.marked(data);
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
<style lang="less">
pre.hljs {
  padding: 12px 2px 12px 40px !important;
  border-radius: 5px !important;
  position: relative;
  font-size: 14px !important;
  line-height: 22px !important;
  overflow: hidden !important;
  code {
    display: block !important;
    margin: 0 10px !important;
    overflow-x: auto !important;
    &::-webkit-scrollbar {
      z-index: 11;
      width: 6px;
    }
    &::-webkit-scrollbar:horizontal {
      height: 6px;
    }
    &::-webkit-scrollbar-thumb {
      border-radius: 5px;
      width: 6px;
      background: #666;
    }
    &::-webkit-scrollbar-corner,
    &::-webkit-scrollbar-track {
      background: #1e1e1e;
    }
    &::-webkit-scrollbar-track-piece {
      background: #1e1e1e;
      width: 6px;
    }
  }
  .line-numbers-rows {
    position: absolute;
    pointer-events: none;
    top: 12px;
    bottom: 12px;
    left: 0;
    font-size: 100%;
    width: 40px;
    text-align: center;
    letter-spacing: -1px;
    border-right: 1px solid rgba(0, 0, 0, 0.66);
    user-select: none;
    counter-reset: linenumber;
    span {
      pointer-events: none;
      display: block;
      counter-increment: linenumber;
      &:before {
        content: counter(linenumber);
        color: #999;
        display: block;
        text-align: center;
      }
    }
  }
  b.name {
    position: absolute;
    top: 2px;
    right: 50px;
    z-index: 10;
    color: #999;
    pointer-events: none;
  }
  .copy-btn {
    position: absolute;
    top: 2px;
    right: 4px;
    z-index: 10;
    color: #333;
    cursor: pointer;
    background-color: #fff;
    border: 0;
    border-radius: 2px;
  }
}
</style>
