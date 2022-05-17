<template>
  <div class="room-container">
    <slot name="room-list-item" v-bind="{ room }">
      <div class="avatar" />
    </slot>
    {{ room.name }}
  </div>
</template>

<script>
export default {
	name: 'RoomsContent',
  props: {
    room: { type: Object, required: true },
  }
}
</script>

<style scoped lang="scss">
.room-container {
	display: flex;
	flex: 1;
	align-items: center;
	width: 100%;
}
.avatar {
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
