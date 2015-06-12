<?php

namespace Persona\Command;


use Persona\Command;
use Persona\Exception\InvalidParametersException;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Filesystem\Filesystem;
use Symfony\Component\Finder\Finder;

class SwitchProfileCommand extends Command
{

    protected function configure()
    {
        $this->setName('switch')
            ->addArgument('profileName', InputArgument::REQUIRED)
            ->setDescription('Switches the settings from files located in the profile name that you type')
        ;

    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $finder = new Finder();
        $fs = new Filesystem();

        $profileName = $input->getArgument('profileName');
        $pathProfiles = $this->get('settings');

        if (!array_key_exists('path_profile', $pathProfiles)) {
            throw new InvalidParametersException(sprintf('%s not found', "path_profile"));
        }

        $fullPath = $pathProfiles['path_profile'].DIRECTORY_SEPARATOR.$profileName;

        $finder->in($fullPath);

        foreach (['files', 'directories'] as $type) {
            $this->createSymlink($fs, $finder, $type);
        }
        $this->updateSettings($fs, $profileName);
    }

    private function createSymlink(Filesystem $fs, Finder $finder, $type = 'files')
    {
        $list = $finder->{$type}();
        foreach ($list as $element) {

            $dotfile = '.'.$element->getRelativePathname();
            $symlinkSetting = $this->get('home').DIRECTORY_SEPARATOR.$dotfile;
            if ($fs->exists($symlinkSetting)) {
                $fs->remove($symlinkSetting);
            }
            $fs->symlink($element->getRealpath(), $symlinkSetting);
        }
    }

    private function updateSettings(Filesystem $fs, $profileName)
    {
        $settings = $this->get('settings');
        $settings['current_profile'] = $profileName;
        $fs->dumpFile($this->get('file_settings'), json_encode($settings));
    }
}