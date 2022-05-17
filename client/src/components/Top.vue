<template>
  <div class="top">
    <div class="rooms-container">
      <h3 style="text-align: center;">
        ルームリスト
      </h3>
      <div class="room-list">
        <div @click="$router.go({path: $router.currentRoute.path, force: true})">
          <room-content :room="{ name: '全体チャット' }" class="room-item" />
        </div>
        <div v-for="room in rooms" :key="room.id" class="room-item" @click="roomChat(room.id)">
          <room-content :room="room" />
        </div>
      </div>
    </div>
    <div class="col-messages">
      <header class="header">
        <h1>Chat</h1>
        <div key="login">
          <button type="button" @click="signOut">
            ログアウト
          </button>
          <button type="button" @click="$router.push('/rooms/new')">
            ルームを作成
          </button>
        </div>
      </header>

      <transition-group name="chat" tag="div" class="list content">
        <section v-for="{ key, name, message } in chat" :key="key" class="item">
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
      <form action="" class="form" @submit.prevent="doSend">
        <textarea v-model="input" :disabled="!user.uid" />
        <button type="submit" :disabled="!user.uid" class="send-button">
          Send
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { getAuth, onAuthStateChanged, signOut } from 'firebase/auth'
import ReconnectingWebSocket from 'reconnecting-websocket';
import RoomContent from './RoomContent.vue'
export default {
  name: 'Top',
  components: {
		RoomContent
	},
  data () {
    return {
      user: {},  // ユーザー情報
      chat: [],  // 取得したメッセージを入れる配列
      rooms: [],
      input: '',  // 入力したメッセージ
      connection: null
    }
  },
  async created() {
    const auth = getAuth();
    onAuthStateChanged(auth, user => {
      this.user = user ? user : {}
      if (user) {
        this.chat = []
      }
    })
    const res = await axios.get(`${process.env.VUE_APP_API_URL}/rooms`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
    })
    this.rooms = res.data

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
    signOut: function () {
      const auth = getAuth();
      signOut(auth).then(() => {
        localStorage.removeItem('jwt')
        this.$router.push('/signin')
      })
    },
    doSend: async function () {
      if (this.user.uid && this.input.length) {
        this.connection.send(this.input);
        this.input = ''
      }
    },
    scrollBottom() {
      this.$nextTick(() => {
        window.scrollTo(0, document.body.clientHeight)
      })
    },
    roomChat(roomId) {
      if (this.connection != null) {
        this.connection.close()
      }

      this.chat = []

      console.log("Starting connection to WebSocket Server")
      const options = { maxRetries: 3 }; // デバッグ時は debug: true を記載する
      this.connection = new ReconnectingWebSocket(
        `${process.env.VUE_APP_WEBSOCKET_URL}/ws/rooms/${roomId}?token=${localStorage.getItem('jwt')}`,
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
    }
  }
}
</script>

<style scoped lang="scss">
.top {
  width: 100%;
  height: 100%;
  display: flex;
  .rooms-container {
    display: flex;
    flex-flow: column;
    flex: 0 0 25%;
    min-width: 260px;
    max-width: 500px;
    position: relative;
    background: #fff;
    height: 100%;

    .room-list {
      flex: 1;
      position: relative;
      max-width: 100%;
      cursor: pointer;
      padding: 0 10px 5px;
      overflow-y: auto;
    }

    .room-item {
      border-radius: 8px;
      align-items: center;
      display: flex;
      flex: 1 1 100%;
      margin-bottom: 5px;
      padding: 0 14px;
      position: relative;
      min-height: 71px;
      transition: background-color 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);

      &:hover {
        background: #f6f6f6;
      }
    }
  }

  .col-messages {
    position: relative;
    height: 100%;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-flow: column;
  }
}
.header {
  width: 100%;
  background: #3ab383;
  margin-bottom: 2em;
  padding: 0.4em 0.8em;
  color: #fff;
}
.form {
  margin-top:auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  bottom: 0;
  height: 80px;
  width: 100%;
  background: #f5f5f5;
}
.form textarea {
  border: 1px solid #ccc;
  border-radius: 2px;
  height: 3em;
  width: 90%;
  resize: none;
  margin-left: 10px;
}
button {
  margin: 10px 10px;
  padding: 10px;
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
/* トランジション用スタイル */
.chat-enter-active {
  transition: all 1s;
}
.chat-enter {
  opacity: 0;
  transform: translateX(-1em);
}
</style>
