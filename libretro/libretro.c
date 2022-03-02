#include "libretro.h"
#include <string.h>
#include <stdio.h>

// callbacks
static retro_video_refresh_t retro_video_refresh = NULL;
static retro_input_poll_t retro_input_poll = NULL;
static retro_input_state_t retro_input_state = NULL;

// input params
static unsigned input_port = 0;
static unsigned input_device = 0;

RETRO_API void retro_set_environment(retro_environment_t cb)
{
    // TODO: implement this
}

RETRO_API void retro_set_video_refresh(retro_video_refresh_t cb)
{
    retro_video_refresh = cb;
}

RETRO_API void retro_set_audio_sample(retro_audio_sample_t cb)
{
    // TODO: implement this
}

RETRO_API void retro_set_audio_sample_batch(retro_audio_sample_batch_t cb)
{
    // TODO: implement this
}

RETRO_API void retro_set_input_poll(retro_input_poll_t cb)
{
    retro_input_poll = cb;
}

RETRO_API void retro_set_input_state(retro_input_state_t cb)
{
    retro_input_state = cb;
}

RETRO_API void retro_init(void)
{
    void Initialize(void);
    Initialize();
}

RETRO_API void retro_deinit(void)
{
    void Deinitialize(void);
    Deinitialize();
}

RETRO_API unsigned retro_api_version(void)
{
    return RETRO_API_VERSION;
}

RETRO_API void retro_get_system_info(struct retro_system_info *info)
{
    memset(info, 0, sizeof(struct retro_system_info));
    void GetEmulatorInfo(struct retro_system_info * info);
    GetEmulatorInfo(info);
}

RETRO_API void retro_get_system_av_info(struct retro_system_av_info *info)
{
    memset(info, 0, sizeof(struct retro_system_av_info));
    void GetEmulatorAVInfo(struct retro_system_av_info * info);
    GetEmulatorAVInfo(info);
}

RETRO_API void retro_set_controller_port_device(unsigned port, unsigned device)
{
    input_port = port;
    input_device = device;
}

RETRO_API void retro_reset(void)
{
    void Reset(void);
    Reset();
}

RETRO_API void retro_run(void)
{
    void Run(void);
    Run();
}

RETRO_API size_t retro_serialize_size(void)
{
    return 0; // TODO: implement this
}

RETRO_API bool retro_serialize(void *data, size_t size)
{
    return true; // TODO: implement this
}

RETRO_API bool retro_unserialize(const void *data, size_t size)
{
    return true; // TODO: implement this
}

RETRO_API void retro_cheat_reset(void)
{
    // TODO: implement this
}

RETRO_API void retro_cheat_set(unsigned index, bool enabled, const char *code)
{
    // TODO: implement this
}

RETRO_API bool retro_load_game(const struct retro_game_info *game)
{
    bool LoadGame(struct retro_game_info *);
    return LoadGame((struct retro_game_info *)game);
}

RETRO_API bool retro_load_game_special(unsigned game_type, const struct retro_game_info *info, size_t num_info)
{
    return true; // TODO: implement this
}

RETRO_API void retro_unload_game(void)
{
    // TODO: implement this
}

RETRO_API unsigned retro_get_region(void)
{
    return RETRO_REGION_NTSC;
}

RETRO_API void *retro_get_memory_data(unsigned id)
{
    return NULL; // TODO: implement this
}

RETRO_API size_t retro_get_memory_size(unsigned id)
{
    return 0; // TODO: implement this
}

void VideoRefresh(const void *data, unsigned width, unsigned height, size_t pitch)
{
    retro_video_refresh(data, width, height, pitch);
}

void InputPoll(void)
{
    retro_input_poll();
}

int16_t InputState(unsigned id)
{
    return retro_input_state(0, RETRO_DEVICE_JOYPAD, 0, id);
}