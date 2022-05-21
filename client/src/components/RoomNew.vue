<template>
  <div class="room-new">
    <img alt="Vue logo" src="../assets/logo.png">
    <h1>グループを作成</h1>
    <button type="button" @click="$router.push('/')">
      トップページ
    </button>
    <ul>
      <input v-model="name" type="text" placeholder="グループ名">
      <li v-for="user in users" :key="user.id">
        {{ user.name }}
        <input v-model="checkedUserIds" type="checkbox" :value="user.id">
      </li>
      <button @click="createRoom">
        作成する
      </button>
    </ul>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'RoomNew',
  data () {
    return {
      users: [],
      name: '',
      checkedUserIds: []
    }
  },
  async created() {
    const res = await axios.get(`${process.env.VUE_APP_API_URL}/users`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
    })
    this.users = res.data
  },
  methods: {
    createRoom: async function() {
      try {
        await axios.post(`${process.env.VUE_APP_API_URL}/rooms`, { name: this.name, userIds: this.checkedUserIds }, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('jwt')}` }
        })
        this.$router.push('/')
      } catch (error) {
          alert(error.response.data.ValidationErrors.join('\n'))
      }
    }
  },
}
</script>

<style scoped>
.room-new {
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
input {
  margin: 10px 0;
  padding: 10px;
}
</style>
