<template>
  <div class="top">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>トップページ</h1>
    <button @click="signOut">
      ログアウト
    </button>
    <button @click="$router.push('/chat')">
      全体チャット
    </button>
    <button type="button" @click="$router.push('/rooms/new')">
      ルームを作成
    </button>
    <h2>ルームリスト</h2>
    <div v-for="room in rooms" :key="room.id" class="room-item" @click="$router.push(`/rooms/${room.id}`)">
      <div class="vac-avatar" />
      {{ room.name }}
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import { getAuth, signOut } from 'firebase/auth'
export default {
  name: 'Top',
  data () {
    return {
      rooms: []
    }
  },
  async created() {
    const res = await axios.get(`${process.env.VUE_APP_API_URL}/rooms`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
    })
    this.rooms = res.data
  },
  methods: {
    signOut: function () {
      const auth = getAuth();
      signOut(auth).then(() => {
        localStorage.removeItem('jwt')
        this.$router.push('/signin')
      })
    },
  }
}
</script>

<style scoped>
.top {
  text-align: center;
  margin-top: 60px;
}
h1, h2 {
  font-weight: normal;
}
button {
  margin: 10px 0;
  padding: 10px;
}
.room-item {
  border-radius: 8px;
  align-items: center;
  display: flex;
  flex: 1 1 100%;
  margin-bottom: 5px;
  padding: 0 400px;
  position: relative;
  min-height: 71px;
  transition: background-color 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
}
.room-item:hover {
  background: #f6f6f6;
}
.vac-avatar {
	background-size: cover;
	background-position: center center;
	background-repeat: no-repeat;
	background-color: #ddd;
	height: 42px;
	width: 42px;
	min-height: 42px;
	min-width: 42px;
	margin-right: 15px;
	border-radius: 50%;
  background-image: url("../assets/default_icon.png")
}
</style>
