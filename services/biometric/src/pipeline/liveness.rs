use image::{DynamicImage, RgbImage};

const LIVENESS_SIZE: u32 = 80;

/// Resize to 80×80 and return BGR channel bytes normalized to [0, 1] as f32 vec.
pub fn preprocess_for_liveness(img: &DynamicImage) -> Vec<f32> {
    let resized = img.resize_exact(LIVENESS_SIZE, LIVENESS_SIZE, image::imageops::FilterType::Triangle);
    let rgb = resized.to_rgb8();
    bgr_normalized(&rgb)
}

fn bgr_normalized(rgb: &RgbImage) -> Vec<f32> {
    let mut out = Vec::with_capacity((rgb.width() * rgb.height() * 3) as usize);
    for pixel in rgb.pixels() {
        out.push(pixel[2] as f32 / 255.0);
        out.push(pixel[1] as f32 / 255.0);
        out.push(pixel[0] as f32 / 255.0);
    }
    out
}

#[cfg(test)]
mod tests {
    use super::*;
    use image::RgbImage;

    #[test]
    fn preprocess_liveness_output_size() {
        let img = DynamicImage::ImageRgb8(RgbImage::new(160, 120));
        let data = preprocess_for_liveness(&img);
        assert_eq!(data.len(), 80 * 80 * 3);
    }
}
