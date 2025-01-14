import React from "react";
import { Button } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Icon } from "../../../../components/Shared/Icon";
import { LoadingIndicator } from "../../../../components/Shared/LoadingIndicator";
import { useToast } from "../../../../hooks/Toast";
import { mutateGenerateVHSClips } from "../../../../core/StashService";

const VHSReelModule = () => {
  const [isGenerating, setIsGenerating] = React.useState(false);
  const Toast = useToast();

  const generateClips = async () => {
    try {
      setIsGenerating(true);
      await mutateGenerateVHSClips();
      Toast.success("Started generating VHS clips");
    } catch (e) {
      Toast.error("Error starting VHS clip generation");
      console.error(e);
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="setting-section">
      <h1>VHS Reel Preparation</h1>
      <div className="setting">
        <div>
          <h3>Generate VHS Clips</h3>
          <div className="sub-heading">
            <FormattedMessage id="Creates 5-second clips for each marker in your library." />
            <FormattedMessage id="Required for VHS Reel playback." />
          </div>
        </div>
        <div>
          <Button
            variant="secondary"
            onClick={() => generateClips()}
            disabled={isGenerating}
          >
            {isGenerating ? (
              <LoadingIndicator message="Generating..." inline small />
            ) : (
              <>
                <Icon icon="tasks" />
                <span><FormattedMessage id="Generate" /></span>
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default VHSReelModule;