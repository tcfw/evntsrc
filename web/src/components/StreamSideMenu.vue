<template>
  <div class="stream-side-menu">
    <div class="search-container">
      <input
        name="search"
        type="text"
        placeholder="Search..."
        v-model="searchInput"
        @keyup.esc="clearSearch"
      />
    </div>
    <div class="stream-list" v-if="!loading">
      <router-link
        :to="'/streams/' + stream.id"
        v-for="stream in filteredStreams"
        :key="stream.id"
        class="stream"
      >
        <div class="icon" :style="iconStyling(stream)">
          <i
            v-if="'icon' in stream && stream.icon == ''"
            :class="'fas fa-' + stream.icon"
          ></i>
          <i v-else class="fas fa-bolt"></i>
        </div>
        <div class="info">
          <div class="name">{{ stream.name }}</div>
          <div class="cluster">{{ stream.cluster }}</div>
        </div>
        <router-link :to="'/streams/' + stream.id + '/settings'"
          ><i class="fas fa-cog config"></i
        ></router-link>
      </router-link>
      <div class="create-button" @click="openCreate"><i class="fas fa-plus"></i></div>
    </div>
    <div class="searching-loader" v-if="loading">
      <i class="fas fa-sync fa-spin"></i>
    </div>
    <div class="loading-error" v-if="error">
      <i class="fas fa-exclamation-triangle" :style="{ color: '#f56c6c' }"></i>
      <p>{{ error }}</p>
      <el-button type="danger" size="mini" @click="retry">Retry</el-button>
    </div>
    <create-modal ref="createModal"></create-modal>
    <key-created-modal :response="keyCreateResponse" @hidden="clearKeyCreateResponse" ref="keyCreatedModal"></key-created-modal>
  </div>
</template>
<script>
import createModal from '@/components/Streams/CreateModal.vue'
import keyCreatedModal from '@/components/Streams/KeyCreatedModal.vue'

export default {
  name: "stream-side-menu",
  components: {
      createModal,
      keyCreatedModal
  },
  data() {
    return {
      searchInput: "",
      loading: false,
      error: null,
      isSearchResults: false,
      streams: [],
      keyCreateResponse: {},
    };
  },
  methods: {
    clearKeyCreateResponse() {
      this.keyCreateResponse = {};
    },
    openCreate() {
      this.$refs.createModal.show();
    },
    clearSearch() {
      this.searchInput = "";
      this.load();
    },
    iconStyling(stream) {
      var styling = {};
      if ("Color" in stream) {
        styling.background = stream.color;
      }
      return styling;
    },
    load() {
      this.loading = true;
      this.error = null;
      this.isSearchResults = false;
      let url = this.$config.API + "/streams";
      if (this.searchInput != "") {
        url += "/search?query=" + this.searchInput;
      }
      this.$http
        .get(url)
        .then(d => {
          this.buildList(d.data);
        })
        .catch(() => {
          this.error =
            "Failed to " +
            (this.searchInput ? "search" : "load") +
            " streams :(";
        })
        .finally(() => {
          if (this.searchInput != "") {
            this.isSearchResults = true;
          }
          this.loading = false;
        });
    },
    buildList(streams) {
      this.streams = streams.streams;
    },
    retry() {
      if (this.searchInput != "") {
        this.search();
      } else {
        this.load();
      }
    }
  },
  computed: {
    filteredStreams() {
      if (this.searchInput == "") {
        return this.streams;
      }
      return this._.filter(
        this.streams,
        stream => stream.name.search(new RegExp(this.searchInput, "gi")) >= 0
      );
    }
  },
  mounted() {
    this.load();
  }
};
</script>

<style lang="scss" scoped>
.stream-side-menu {
  @apply bg-white w-full h-full shadow-lg relative float-left;
}

.create-button {
  @apply mx-4 mt-8 p-2 bg-ev-100 rounded text-white block text-center text-sm cursor-pointer;

  &:hover {
    @apply bg-ev-50;
  }
}

.search-container {
  @apply absolute shadow-inner rounded border border-gray-200 m-4;
  width: calc(100% - 30px);
  height: 36px;
  background: #fafafa;
  z-index: 2;

  input[type="text"] {
    background: none;
    border: none;
    color: #334;
    height: calc(100% - 2px);
    margin-left: 10px;
    outline: none;
    font-weight: 400;
    font-size: 12px;

    &::placeholder {
      font-weight: 100;
      color: #bbc;
    }
  }
  .fas {
    @apply absolute text-ev-100 cursor-pointer;
    right: 12px;
    top: 10px;
    font-size: 14px;
  }

  .clear-search {
    right: 35px;
  }
}

.searching-loader {
  @apply absolute text-ev-100;
  position: absolute;
  left: 50%;
  top: 80px;
  transform: translateX(-10px);
  font-size: 14px;
}

.loading-error {
  color: white;
  font-size: 12px;
  position: absolute;
  top: 80px;
  text-align: center;
  width: 100%;
}

.stream-list {
  padding-top: 70px;
  position: relative;
  width: 100%;

  .stream {
    @apply relative cursor-pointer block w-full;
    padding: 10px 15px;

    &:hover,
    &.router-link-active {
      @apply bg-ev-900 border border-0 border-r-4 border-ev-700;

      .fas.config {
        right: 11px;
      }
    }

    .icon,
    .info {
      display: inline-block;
      vertical-align: middle;
    }

    .icon {
      @apply bg-ev-40 rounded;
      height: 42px;
      width: 42px;

      .fas {
        @apply text-ev-900 w-full h-full text-center;
        line-height: 42px;
        font-size: 16px;
      }
    }

    .info {
      margin-left: 10px;

      .name {
        @apply text-text-200;
        font-size: 12px;
        font-weight: 300;
      }

      .cluster {
        @apply text-text-400 uppercase;
        font-size: 9px;
        margin-top: -1px;
      }
    }

    .fas.config {
      @apply text-ev-100 absolute cursor-pointer;
      top: 50%;
      transform: translateY(-50%);
      right: 15px;

      &:hover {
        @apply text-white;
      }
    }
  }
}
</style>
