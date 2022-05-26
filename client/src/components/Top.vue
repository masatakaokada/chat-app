<template>
  <v-row class="fill-height" no-gutters>
    <v-col cols="12" sm="3" class="border">
      <v-container>
        <h3 style="text-align: center;">
          Room List
        </h3>
        <v-list two-line>
          <v-list-item
            :prepend-avatar="avatar"
            title="全体チャット"
            @click="$router.go({path: $router.currentRoute.path, force: true})"
          />
          <v-divider />
          <div v-for="(room, index) in rooms" :key="room.id">
            <v-list-item
              :prepend-avatar="avatar"
              :title="room.name"
              @click="roomChat(room.id, room.name)"
            >
              <!-- <v-list-item-icon icon="mdi-dots-vertical" /> -->
            </v-list-item>
            <v-divider v-if="index < rooms.length - 1" :key="index" />
          </div>
        </v-list>
      </v-container>
    </v-col>
    <v-col cols="12" sm="9" class="d-flex flex-column">
      <!-- ヘッダー -->
      <v-toolbar color="#3ab383">
        <v-toolbar-title class="text-white">
          {{ currentRoomName }}
        </v-toolbar-title>
        <v-btn class="text-white" type="button" @click="$router.push('/rooms/new')">
          ルームを作成
        </v-btn>
        <v-btn class="text-white" type="button" @click="signOut">
          ログアウト
        </v-btn>
      </v-toolbar>

      <!-- チャットメッセージ -->
      <v-container id="message-area" class="flex-grow-1 overflow-y-auto" style="height: 80vh;">
        <v-list lines="three">
          <div v-for="{ key, id, message } in chat" :key="key">
            <v-list-item :prepend-avatar="id == currentUser.id ? null : avatar">
              <v-card
                :class="id == currentUser.id ? 'ml-auto pa-5' : 'pa-5'"
                :color="id == currentUser.id ? '#B9F6CA' : '#F5F5F5'"
                max-width="350px"
                flat
              >
                <v-list-item-title style="white-space:unset;">
                  {{ message }}
                </v-list-item-title>
              </v-card>
            </v-list-item>
          </div>
        </v-list>
      </v-container>

      <!-- 入力フォーム -->
      <!-- <v-container class="mt-auto"> -->
      <v-container>
        <v-text-field
          v-model="input"
          hide-details
          append-icon="mdi-send"
          label="Message"
          type="text"
          @click:append="doSend"
        />
      </v-container>
    </v-col>
  </v-row>
</template>

<script>
import axios from 'axios'
import { getAuth, onAuthStateChanged, signOut } from 'firebase/auth'
import ReconnectingWebSocket from 'reconnecting-websocket';
export default {
  name: 'Top',
  data () {
    return {
      user: {},  // ユーザー情報
      chat: [],  // 取得したメッセージを入れる配列
      rooms: [],
      input: '',  // 入力したメッセージ
      connection: null,
      avatar: require("@/assets/default_icon.png"),
      currentRoomName: '全体チャット',
      currentUser: {}
    }
  },
  async created() {
    const auth = getAuth();
    onAuthStateChanged(auth, user => {
      this.user = user ? user : {}
    })
    if (this.user) {
      this.chat = []
      const res = await axios.get(`${process.env.VUE_APP_API_URL}/user`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
      })
      this.currentUser = res.data
    }
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
        id: obj.id,
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
        const scrollHeight = document.getElementById('message-area').scrollHeight;
        document.getElementById('message-area').scrollTop = scrollHeight;
      })
    },
    roomChat(roomId, roomName) {
      this.currentRoomName = roomName

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
          id: obj.id,
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
