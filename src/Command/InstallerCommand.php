<?php

namespace Persona\Command;

use Persona\Command;
use Persona\Json\JsonFile;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\Question;
use Symfony\Component\Filesystem\Filesystem;

class InstallerCommand extends Command
{
    protected function configure()
    {
        $this->setName('install')
            ->setDescription('Install the settings file for switching profiles');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $fs = new Filesystem();
        $helper = $this->getHelper('question');
        $baseSettings = $this->get('directory_settings');
        $distSettings = json_decode(file_get_contents($this->get('settings_dist')), true);

        $userSettings = [];
        foreach ($distSettings as $key => $setting) {
            $question = new Question($setting);
            $response = $helper->ask($input, $output, $question);

            if (!$fs->exists($baseSettings)) {
                $fs->mkdir($baseSettings);
                $fs->touch($this->get('file_settings'));
            }

            $userSettings[$key] = $response;
        }

        $updateSettings = new JsonFile($this->get('file_settings'));
        $updateSettings->write($userSettings);
    }

}