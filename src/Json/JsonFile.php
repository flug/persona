<?php


namespace Persona\Json;


use Seld\JsonLint\JsonParser;
use Seld\JsonLint\ParsingException;
use Symfony\Component\Filesystem\Exception\IOException;
use Symfony\Component\Filesystem\Filesystem;

class JsonFile
{

    const LAX_SCHEMA = 1;
    const STRICT_SCHEMA = 2;
    const JSON_UNESCAPED_SLASHES = 64;
    const JSON_PRETTY_PRINT = 128;
    const JSON_UNESCAPED_UNICODE = 256;


    private $fs;
    private $path;

    public function __construct($path)
    {
        $this->path = $path;
        $this->fs = new Filesystem();
    }

    public function read()
    {
        try {
            $json = file_get_contents($this->path);
        } catch (\Exception $e) {
            throw new \RuntimeException('Could not read '.$this->path."\n\n".$e->getMessage());
        }

        return static::parseJson($json, $this->path);
    }

    /**
     * Parses json string and returns hash.
     *
     * @param string $json json string
     * @param string $file the json file
     *
     * @return mixed
     */
    public static function parseJson($json, $file = null)
    {
        if (null === $json) {
            return;
        }
        $data = json_decode($json, true);
        if (null === $data && JSON_ERROR_NONE !== json_last_error()) {
            self::validateSyntax($json, $file);
        }

        return $data;
    }

    /**
     * Validates the syntax of a JSON string
     *
     * @param  string $json
     * @param  string $file
     * @return bool                      true on success
     * @throws \UnexpectedValueException
     * @throws JsonValidationException
     * @throws ParsingException
     */
    protected static function validateSyntax($json, $file = null)
    {
        $parser = new JsonParser();
        $result = $parser->lint($json);
        if (null === $result) {
            if (defined('JSON_ERROR_UTF8') && JSON_ERROR_UTF8 === json_last_error()) {
                throw new \UnexpectedValueException('"'.$file.'" is not UTF-8, could not parse as JSON');
            }

            return true;
        }
        throw new ParsingException('"'.$file.'" does not contain valid JSON'."\n".$result->getMessage(),
            $result->getDetails());
    }

    public function write($data, $options = 448)
    {
        try{
            $this->fs->dumpFile($this->path, static::encode($data, $options). ($options & self::JSON_PRETTY_PRINT ? "\n" : ''));
        }catch (\Exception $e)
        {
            throw new IOException(sprintf('Failed to write file "%s".', $this->path), 0, null, $this->path);
        }
    }

    public static function encode($data, $options = 448)
    {

        $json = json_encode($data);
        if (false === $json) {
            self::throwEncodeError(json_last_error());
        }

        $prettyPrint = (bool)($options & self::JSON_PRETTY_PRINT);
        $unescapeUnicode = (bool)($options & self::JSON_UNESCAPED_UNICODE);
        $unescapeSlashes = (bool)($options & self::JSON_UNESCAPED_SLASHES);
        if (!$prettyPrint && !$unescapeUnicode && !$unescapeSlashes) {
            return $json;
        }

        $result = JsonFormatter::format($json, $unescapeUnicode, $unescapeSlashes);
        return $result;

    }

    /**
     * Throws an exception according to a given code with a customized message
     *
     * @param int $code return code of json_last_error function
     * @throws \RuntimeException
     */
    private static function throwEncodeError($code)
    {
        switch ($code) {
            case JSON_ERROR_DEPTH:
                $msg = 'Maximum stack depth exceeded';
                break;
            case JSON_ERROR_STATE_MISMATCH:
                $msg = 'Underflow or the modes mismatch';
                break;
            case JSON_ERROR_CTRL_CHAR:
                $msg = 'Unexpected control character found';
                break;
            case JSON_ERROR_UTF8:
                $msg = 'Malformed UTF-8 characters, possibly incorrectly encoded';
                break;
            default:
                $msg = 'Unknown error';
        }
        throw new \RuntimeException('JSON encoding failed: '.$msg);
    }

    public function exists()
    {
        return $this->fs->exists($this->path);
    }

    public function getPath()
    {
        return $this->path;
    }
}