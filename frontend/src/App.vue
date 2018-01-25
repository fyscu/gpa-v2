<!-- 首页内容渲染 -->
<template>
  <div>
  
    <div id="header">
      <appBar :zDepth="0" title="四川大学成绩/绩点一键计算" class="appbar">
        <iconButton @click="logout" icon=':icon-logout' iconClass="iconfont" slot="right"></iconButton>
      </appBar>
  
    </div>
  
    <!-- 渲染出口 -->
    <div class="content">
      <div class="wrapper">
        <router-view></router-view>
  
      </div>
    </div>
  
  </div>
</template>
<script>
// import './css/iconfont.css'
import appBar from "muse-components/appBar"
import iconButton from "muse-components/iconButton"

export default {
  components: {
    appBar,
    iconButton
  },
  data() {
    const desktop = isDesktop();
    return {
      open: desktop,
      docked: desktop,
      desktop: desktop,
    }
  },
  methods: {
    logout() {
      localStorage.clear()
      location.reload()
    }
  },
  mounted() {
    this.routes = this.$router.options.routes
    window.addEventListener('resize', this.handleResize)
    window.addEventListener('hashchange', () => {
    })
  }

}

function isDesktop() {
  return window.innerWidth > 993
}

</script>
<style lang="less">
.content {
  transition: all .45s @easeOutFunction;
  &.nav-hide {
    padding-left: 0;
  }
}

.wrapper {
  padding: 48px 72px;
}

@media (min-width: 480px) {
  .content {
    max-width: 1200px;
    margin: 0 auto;
  }
}

@media (max-width: 993px) {
  .appbar {
    left: 0;
  }
  .content {
    padding-left: 0;
  }
  .wrapper {
    padding: 24px 0px;
  }
}

@easeOutFunction: cubic-bezier(0.23, 1, 0.32, 1);
@easeInOutFunction: cubic-bezier(0.445, 0.05, 0.55, 0.95);
</style>
</style>