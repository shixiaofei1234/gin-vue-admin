<template>
  <div class="ai-chat-disabled-page">
    <div class="disabled-card">
      <div class="title">{{ AI_DISABLED_PAGE.title }}</div>
      <div class="desc">
        {{ disabledLines[0] }}
        <strong>{{ AI_ASSISTANT.floatButtonChar }}</strong>
        {{ disabledLines[1] }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { AI_ASSISTANT, AI_DISABLED_PAGE } from '@/constants/aiFloatChat'

defineOptions({
  name: 'AIChat'
})

const disabledLines = computed(() =>
  AI_DISABLED_PAGE.templateLines(AI_ASSISTANT.name, AI_ASSISTANT.floatButtonChar)
)
</script>

<style scoped>
.ai-chat-disabled-page {
  min-height: 360px;
  padding: 16px;
}

.disabled-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
}

.desc {
  color: #4b5563;
  line-height: 1.7;
}
</style>
