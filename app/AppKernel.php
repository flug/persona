<?php

class AppKernel
{

    /**
     * @var array
     */
    public static $defaultConfig = [

        'directory_settings' => '{$home}/.switcher',
        'file_settings' => '{$directory_settings}/switcher.json',
        'profiles_directory' => '{$home}/profiles',
        'settings_dist' => '{$rootKernel}/settings.json.dist',
        'version' => '1.1.0@alpha',
        'name' => 'Persona'

    ];
    private $config;

    public function __construct()
    {
        $this->config = static::$defaultConfig;
    }

    /**
     * @return static
     */
    public static function getInstance()
    {
        return new static();
    }

    /**
     * @return array of Command
     */
    public function loadCommands()
    {
        $commands = [
            new Persona\Command\InstallerCommand(),
        ];
        if ($this->settingsFileExist()) {
            $commands = array_merge([
                    new Persona\Command\SwitchProfileCommand(),
                    new Persona\Command\AddProfileCommand(),
                ], $commands);

        }

        return $commands;
    }

    private function settingsFileExist()
    {
        if (!(new \Persona\Json\JsonFile($this->get('file_settings')))->exists()) {
            return false;
        }
        return true;
    }

    /**
     * @param $key
     * @param int $flags
     * @return mixed|string
     */
    public function get($key, $flags = 0)
    {

        switch ($key) {
            case 'home':
                return rtrim(getenv('HOME') ?: getenv('USERPROFILE'), '/\\');
                break;
            case 'rootKernel':
                return __DIR__;
                break;
            case 'settings':
                return (new \Persona\Json\JsonFile($this->get('file_settings')))->read();
                break;
            default:
            case 'directory_settings':
            case 'profiles_directory':
            case 'file_settings':
                return $this->process($this->config[$key], $flags);
                break;
        }
    }

    /**
     * @param $value
     * @param $flags
     * @return mixed
     */
    private function process($value, $flags)
    {
        $config = $this;
        if (!is_string($value)) {
            return $value;
        }

        return preg_replace_callback('#\{\$(.+)\}#', function ($match) use ($config, $flags) {
            return $config->get($match[1], $flags);
        }, $value);
    }
}