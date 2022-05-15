<template>
  <div>
    <header class="header">
      <h1>Chat</h1>
      <div
        v-if="user.uid"
        key="login"
      >
        [{{ user.displayName ?? '名無しさん' }}]
        <button
          type="button"
          @click="doLogout"
        >
          ログアウト
        </button>
        <button
          type="button"
          @click="$router.push('/')"
        >
          トップページ
        </button>
      </div>
    </header>

    <transition-group
      name="chat"
      tag="div"
      class="list content"
    >
      <section
        v-for="{ key, name, message } in chat"
        :key="key"
        class="item"
      >
        <div class="item-image">
          <img src="../assets/default_icon.png" width="40" height="40">
        </div>
        <div class="item-detail">
          <div class="item-name">
            {{ name }}
          </div>
          <div class="item-message">
            {{ message }}
          </div>
        </div>
      </section>
    </transition-group>
  
    <!-- 入力フォーム -->
    <form
      action=""
      class="form"
      @submit.prevent="doSend"
    >
      <textarea
        v-model="input"
        :disabled="!user.uid"
      />
      <button
        type="submit"
        :disabled="!user.uid"
        class="send-button"
      >
        Send
      </button>
    </form>
  </div>
</template>

<script>
import { getAuth, onAuthStateChanged, signOut } from "firebase/auth";
import ReconnectingWebSocket from 'reconnecting-websocket';
export default {
  data() {
    return {
      user: {},  // ユーザー情報
      chat: [],  // 取得したメッセージを入れる配列
      input: '',  // 入力したメッセージ
      connection: null
    }
  },
  created() {
    const auth = getAuth();
    onAuthStateChanged(auth, user => {
      this.user = user ? user : {}
      if (user) {
        this.chat = []
      }
    })

    console.log("Starting connection to WebSocket Server")
    const options = { maxRetries: 3 }; // デバッグ時は debug: true を記載する
    this.connection = new ReconnectingWebSocket(
      `${process.env.VUE_APP_WEBSOCKET_URL}/ws?token=${localStorage.getItem('jwt')}`,
      [],
      options
    )

    this.connection.onmessage = event => {
      const obj = JSON.parse(event.data);
      this.chat.push({
        key: this.$uuid.v4(),
        name: obj.username,
        message: obj.message
      })
      this.scrollBottom()
    }

    this.connection.onopen = () => {
      console.log("Successfully connected to the websocket server")
    }

    this.connection.onclose = () => {
      console.log("connection closed")
    }

    this.connection.onerror = error => {
      console.log("there was an error")
      console.log(error)
    }
  },
  methods: {
    doLogout: function () {
      const auth = getAuth();
      signOut(auth).then(() => {
        localStorage.removeItem('jwt')
        this.$router.push('/signin')
      })
    },
    scrollBottom() {
      this.$nextTick(() => {
        window.scrollTo(0, document.body.clientHeight)
      })
    },
    doSend: async function () {
      if (this.user.uid && this.input.length) {
        this.connection.send(this.input);
        this.input = ''
      }
    },
  }
}
</script>

<style scoped>
* {
  margin: 0;
  box-sizing: border-box;
}
.header {
  background: #3ab383;
  margin-bottom: 2em;
  padding: 0.4em 0.8em;
  color: #fff;
}
.content {
  margin: 0 auto;
  padding: 0 10px;
  max-width: 600px;
}
.form {
  position: fixed;
  display: flex;
  justify-content: center;
  align-items: center;
  bottom: 0;
  height: 80px;
  width: 100%;
  background: #f5f5f5;
}
.form textarea {
  border: 1px solid #ccc;
  border-radius: 2px;
  height: 4em;
  width: calc(100% - 6em);
  resize: none;
}
.list {
  margin-bottom: 100px;
}
.item {
  position: relative;
  display: flex;
  align-items: flex-end;
  margin-bottom: 2em;
}
.item-image img {
  border-radius: 20px;
  vertical-align: top;
}
.item-detail {
  margin: 0 0 0 1.4em;
}
.item-name {
  font-size: 75%;
}
.item-message {
  position: relative;
  display: inline-block;
  padding: 0.8em;
  background: #deefe8;
  border-radius: 4px;
  line-height: 1.2em;
}
.item-message::before {
  position: absolute;
  content: " ";
  display: block;
  left: -16px;
  bottom: 12px;
  border: 4px solid transparent;
  border-right: 12px solid #deefe8;
}
.send-button {
  height: 4em;
}

button {
  margin: 10px 0;
  padding: 10px;
}
/* トランジション用スタイル */
.chat-enter-active {
  transition: all 1s;
}
.chat-enter {
  opacity: 0;
  transform: translateX(-1em);
}
</style>
